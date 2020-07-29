#!/usr/bin/perl
use encoding 'utf8', STDIN=>'utf8', STDOUT=>'utf8';

use Data::Dumper;
use File::Slurp;
use JSON::XS;
use POSIX;

use Digest::SHA qw(sha1 sha1_hex sha1_base64);
use Compress::Zlib;
use MIME::Base64;
use Encode;
use List::Util qw(min max sum);

my $left_jss	= shift;
my $options		= shift;
my $line_offset	= shift || 0;
my $line_cutoff	= shift || 0;

our $cached 	= 1;
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

	print STDERR "Parsing $filename\n";

	open(my $handle, '<', $filename) or die($!);
	binmode($handle,":encoding(utf8)");
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

our $processor = {};

sub process_idfa
{
	my $entry = shift;
	my $value = shift;
	my $json = shift;

	if($json->{platform} eq "android")
	{
		if($json->{googleadid} and $json->{googleadid} ne "0" and $json->{googleadid} ne "00000000-0000-0000-0000-000000000000")
		{
			$json->{udid} = $json->{googleadid};
			$json->{udid_t} = "gaid";
			#print STDERR $json->{requestid}." udid-gaid ".$json->{udid}."\n";
		}
		elsif($json->{imei} and $json->{imei} ne "0")
		{
			$json->{udid} = $json->{imei};
			$json->{udid_t} = "imei";
			#print STDERR $json->{requestid}." udid-imei ".$json->{udid}."\n";
		}
		elsif($json->{androidid} and $json->{androidid} ne "0" and $json->{androidid} ne "00000000-0000-0000-0000-000000000000")
		{
			$json->{udid} = $json->{androidid};
			$json->{udid_t} = "aid";
			#print STDERR $json->{requestid}." udid-androidid ".$json->{udid}."\n";
		}
		elsif($json->{sysid} and $json->{sysid} ne ",")
		{
			$json->{udid} = $json->{sysid};
			$json->{udid_t} = "a_sysid";
			#print STDERR $json->{requestid}." udid-asysid ".$json->{udid}."\n";
		}
		else
		{
			$json->{udid} = "";
			$json->{udid_t} = "a_ipua";
			#print STDERR $json->{requestid}." ANDROID-IPUA\n";
		}
	}
	else
	{
		if($value ne "0" and $value ne "" and $value ne "00000000-0000-0000-0000-000000000000")
		{
			$json->{udid} = $value;
			$json->{udid_t} = "idfa";
		}
		elsif($json->{sysid} and $json->{sysid} ne ",")
		{
			$json->{udid} = $json->{sysid};
			$json->{udid_t} = "i_sysid";
			#print STDERR $json->{requestid}." udid-isysid ".$json->{udid}."\n";
		}
		else
		{
			$json->{udid} = "";
			$json->{udid_t} = "i_ipua";
			#print STDERR $json->{requestid}." IOS-IPUA\n";
		}
	}

}
$processor->{idfa} = \&process_idfa;

sub process_requestid
{
	my $entry = shift;
	my $value = shift;
	my $json = shift;

	my $v2 = $json->{campaignid};

	my $v = "$value:$v2";
	my $id = sha1_hex(encode_base64(encode_utf8($v), ''));

	$json->{id} = $id;
}
$processor->{requestid} = \&process_requestid;

sub process_cc
{
	my $entry = shift;
	my $value = shift;
	my $json = shift;

	$json->{$entry} = lc $value;
	#print STDERR "$entry -> $value :: ".$json->{$entry}."\n";
}
$processor->{countrycode} = \&process_cc;

sub process_locale
{
	my $entry = shift;
	my $value = shift;
	my $json = shift;

	my @locales = split("-", $value);

	$json->{locale} = lc $locales[0];
	$json->{$entry} = lc $value;
	#print STDERR "locale -> $value :: ".$json->{locale}."\n";
}
$processor->{language} = \&process_locale;

sub process_osversion
{
	my $entry = shift;
	my $value = shift;
	my $json = shift;

	my @versions = split("[.]", $value);
	$json->{osv} = $json->{platform}.$versions[0];
	$json->{$entry} = $json->{platform}.$value;
	#print STDERR "osv -> $value :: ".$json->{osv}."\n";
	#print STDERR "$entry -> $value :: ".$json->{$entry}."\n";
}
$processor->{osversion} = \&process_osversion;

sub process_screensize
{
	my $entry = shift;
	my $value = shift;
	my $json = shift;

	my @sizes = split("x", $value);

	if(scalar @sizes == 2)
	{
		my $w = $sizes[0];
		my $h = $sizes[1];

		if($w >= $h)
		{
			$json->{scr_o} = "H_W";
			$json->{scr_w} = $h;
		}
		else
		{
			$json->{scr_o} = "W_H";
			$json->{scr_w} = $w;
		}

	}
	else
	{
		$json->{scr_o} = "H_W";
		$json->{scr_w} = $w;
	}
	#print STDERR "scr_o -> $value :: ".$json->{scr_o}."\n";
	#print STDERR "scr_w -> $value :: ".$json->{scr_w}."\n";
}
$processor->{screensize} = \&process_screensize;


sub parseLeftRow
{
	my $line = shift;
	#chomp($line);

	my $json = $JSON_READER->decode($line);

	for my $entry (keys %$json)
	{
		if(exists $processor->{$entry})
		{
			$processor->{$entry}->($entry, $json->{$entry}, $json);
		}
	}

	print $JSON_LOGGER->encode($json)."\n";
}


sub ApplyConfig
{
	my $config = shift;

	if(exists $config->{"verbose"})
	{
		$verbose = $config->{"verbose"};
	}

	if(exists $config->{"cached"})
	{
		$cached = $config->{"cached"};
	}

	print STDERR "config => \"verbose:$verbose;cached:$cached;\"\n";
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