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

my $codec = {
	metadata=>["rts", "day", "algorithm", "deviceip", "requestid", "udid", "udid_t", "ivr_model", "time"],
	features=>[	"adtype", "advertiserid", "appid",
				"campaignid", "carrier", "citycode", "countrycode", "creativeid",
				"devicetype",
				"exchanges", "language", "locale",
				"make", "model",
				"networktype",
				"osv", "osversion",
				"packagename", "placementid", "platform", "publisherid",
				"scr_o", "scr_w", "screensize", "sdkversion", "size",
				"templateid", "traffictype",
				"udid_t", "unitid"
				]
};

our $wasted_verbose = {};

sub parseLeftRow
{
	my $line = shift;
	#chomp($line);
	my $json = $JSON_READER->decode($line);

	my $feature_ = [];
	my $metadata_ = {};
	my $recoded = { id=>$json->{id}, target=>$json->{target}, variant=>$json->{variant}, meta=>$metadata_};

	my $locus = 0;
	foreach my $feature (@{$codec->{features}})
	{
		# SEMANTIC : [FEATURE NAME, FEATURE VALUE, FEATURE LOCUS (# of feature), NONE TRIVIAL FEATURE]
		push @$feature_, [ $feature, $json->{$feature}, $locus, $json->{$feature} ? 1 : 0];
		$locus++;
	}

	$digest = sha1_hex( $JSON_LOGGER->encode($feature_) );
	$recoded->{digest} = $digest;
	$recoded->{feature} = $feature_;

	foreach my $metadata (@{$codec->{metadata}})
	{
		$metadata_->{$metadata} = $json->{$metadata};
	}

	$metadata_->{features} = $locus;

	if($verbose)
	{
		foreach my $feature (@{$codec->{features}})
		{
			delete $json->{$feature};
		}

		delete $json->{id};
		delete $json->{rts};
		delete $json->{target};
		delete $json->{requestid};
		delete $json->{variant};

		foreach my $entry (keys %$json)
		{
			$wasted_verbose->{$entry} ++;
		}
	}

	print $JSON_LOGGER->encode($recoded)."\n";
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

	if($verbose)
	{
		print STDERR $JSON_LOGGER->pretty->canonical->encode($wasted_verbose);
	}
}

Main();

1
