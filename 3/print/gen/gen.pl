#!/usr/bin/env perl
# 2014-09-04 Adam Bryt

use strict;
use warnings;

if ($#ARGV != 1) {
	printf(STDERR "usage: gen.pl: NUM STRING\n");
	exit(1);
}

my ($n, $str) = @ARGV;

for my $i (1 .. $n) {
	printf("%6d %s\n", $i, $str);
}
