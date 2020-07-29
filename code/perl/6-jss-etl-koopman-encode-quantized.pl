#!/usr/bin/perl
use encoding 'utf8', STDIN=>'utf8', STDOUT=>'utf8';

use Data::Dumper;
use File::Slurp;
use JSON::XS;
use POSIX;
use IO::Zlib;

use Digest::SHA qw(sha1 sha1_hex sha1_base64);
use Compress::Zlib;
use MIME::Base64;
use Encode;
use List::Util qw(min max sum);

my $left_jss	= shift;
my $options		= shift;
my $line_offset	= shift || 0;
my $line_cutoff	= shift || 0;

our $cardinal	= 24;
our $bandwidth	= 64;
our $quantized 	= 256;
our $barrier 	= 32;
our $cached 	= 1;
our $gzip = 0;
our $verbose 	= 0;

our $JSON_READER = JSON::XS->new;
our $JSON_LOGGER = JSON::XS->new->canonical->utf8;

# ++++++++++ Synthetic Option API ++++++++++

sub parseSynthetic
{
	my $option_line = shift;

	return {} unless $option_line;

	if($option_line =~ m/\:/)
	{
		my @options = split(";", $option_line);
		my $result = {};
		%$result = map { my @pair = split(":", $_); $pair[0] => $pair[1] } @options;
		return $result;
	}

	my $text = read_file($option_line, binmode => ':raw');

	die "cannot read file $option_line" unless $text;

	return JSON::XS->new->allow_nonref->utf8->decode($text);
}

# ---------- Synthetic Option API ----------

sub loadJSON
{
	my $file = shift || die "crap";

	my $text = read_file($file, binmode => ':raw');

	die "cannot read file $file" unless $text;

	return JSON::XS->new->allow_nonref->utf8->decode($text);
}

sub encodeLine
{
	my $json = shift;

	my $line2 = $JSON_LOGGER->encode($json);
	my $code = compress(encode_utf8($line2), 9);
	my $line3 = encode_base64($code);
	$line3 =~ s/\n//g;

	return $line3;
}

sub decodeLine
{
	my $line = shift;

	my $code1 = decode_base64($line);
	my $code2 = uncompress($code1);
	my $text = decode_utf8($code2);

	return $JSON_PARSER->decode($text);
}

sub parseFile
{
	my $parser = shift;
	my $filename = shift;
	my $offset = shift;
	my $lines = shift || 1024*1024*1024;

	$filename or die "no file to parse";

	my $lineno = 0;
	my $parsed = 0;
	my $handle = undef;	

	if($gzip)
	{
		print STDERR "Parsing $filename with GZIP\n";
		$handle = IO::Zlib->new($filename, "rb");
	}
	else
	{
		print STDERR "Parsing $filename with plain UTF8 lines\n";
		open($handle, '<', $filename) or die($!);
		binmode($handle,":encoding(utf8)");
	}

	while(my $line = <$handle>)
	{
		if($lineno >= $offset)
		{
			$parser->($line, $lineno);
			$parsed ++;
		}
		if($parsed >= $lines)
		{
			$lineno++;
			last;
		}
		if($lineno % 1024 == 1)
		{
			print STDERR ".";
		}
		$lineno++;
	}
	close($handle);

	print STDERR "\n$parsed / $lineno lines parsed\n";

}

our $all_feature = {};

sub encodeSDR
{
	my $signature = shift;
	my $N = shift;

	# READ CACHE
	return $all_feature->{$signature} if $cached and exists $all_feature->{$signature};

	my $vector = [[],[]];
	my $MAX = 1<<32;
	my $DR = ($quantized - $barrier);
	my $occupied = {};
	my $distinctive = 0;

	for(my $i = 0; $i < $N*2 && $distinctive<$N; $i++)
	{
		my $digest = sha1_hex("X".$signature."X".$i."X");
		
		my $hex = substr($digest,16,8);
		my $x = hex($hex)/$MAX;
		$x = floor($bandwidth * $x) % $bandwidth;

		if(exists $occupied->{$x})
		{
			print STDERR "collision #$x => ".$occupied->{$x}."\n" if $verbose;
			next;
		}

		$distinctive++;

		push @{$vector->[0]}, $x;

		my $hexy = substr($digest,32,8);
		my $y = hex($hexy)/$MAX;
		$y = floor($DR * $y) % $DR;
		$y = $y + $barrier;

		if($verbose && $y >= $quantized)
		{
			die "!!! $y > $quantized overflow\n";
		}

		if($verbose && $y == $barrier)
		{
			print STDERR "!!! $y == 32 ground state\n";
		}

		push @{$vector->[1]}, $y;
		$occupied->{$x} = $y;
	}

	if($verbose && scalar @{$vector->[0]} != $N or scalar @{$vector->[0]} != scalar @{$vector->[1]})
	{
		die "Partial SDR with valid = ".(scalar @{$vector->[0]})."\n";
	}

	# WRITE CACHE
	$all_feature->{$signature} = $vector if $cached;

	return $vector;
}

sub parseLeftRow
{
	my $line = shift;
	#chomp($line);

	my $source = $JSON_READER->decode($line);
	my $encoded = {};
	my $sdr = [];

	#this is very very expensive to do

	foreach my $feature (@{$source->{feature}})
	{
		my $pairwise = $feature->[0].",;|".$feature->[1];
		my $precode =  encode_base64(encode_utf8( $pairwise ), '');
		my $v = encodeSDR($precode, $cardinal);
		#print Dumper $v;

		my $locus = $feature->[2];
		if(defined $sdr->[$locus])
		{
			die "multiple value locus not supported yet!";
		}

		$sdr->[$locus]  = $v;
	}

	my $verbatim = scalar @$sdr;

	$encoded->{id} = $source->{id};
	$encoded->{digest} = $source->{digest};
	$encoded->{target} = $source->{target};

	$encoded->{sdr} = $sdr;
	$encoded->{meta} = { cardinal=>$cardinal, bandwidth=>$bandwidth, features=>$source->{meta}->{features}, verbatim=>$verbatim};

	print $JSON_LOGGER->encode($encoded)."\n";	
}

sub ApplyConfig
{
	my $config = shift;

	if(exists $config->{"verbose"})
	{
		$verbose = $config->{"verbose"};
	}

	if(exists $config->{"cardinal"})
	{
		$cardinal = $config->{"cardinal"};
	}

	if(exists $config->{"bandwidth"})
	{
		$bandwidth = $config->{"bandwidth"};
	}

	if(exists $config->{"quantized"})
	{
		$quantized = $config->{"barquantizedrier"};
	}

	if(exists $config->{"barrier"})
	{
		$barrier = $config->{"barrier"};
	}

	if(exists $config->{"cached"})
	{
		$cached = $config->{"cached"};
	}

	if(exists $config->{"gzip"})
	{
		$gzip = $config->{"gzip"};
	}

	print STDERR "config => \"verbose:$verbose;gzip:$gzip;cached:$cached;barrier:$barrier;quantized:$quantized;bandwidth:$bandwidth;cardinal:$cardinal\"\n";
}


sub Main 
{
	$variant = parseSynthetic($options);
	ApplyConfig($variant);

	my $t1 = time();
	parseFile(\&parseLeftRow,$left_jss, $line_offset, $line_cutoff);
	print STDERR "cost = ".(time() - $t1)."s\n";
}

Main();

1