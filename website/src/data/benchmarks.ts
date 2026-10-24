// Results behind the benchmark charts and tables, updated together. The suite reruns with the commands on the
// Benchmarks page; the comparison comes from a harness that runs the old and new versions side by side.

export const environment = 'Apple M4 Max, Go 1.27.0, median of 5 runs';

/** The two versions the charts compare. */
export const versions = { before: 'v1', after: 'v2' } as const;

export type Group = 'location' | 'errors' | 'logging';

export interface Result {
	ns: number;
	allocs: number;
}

export interface Comparison {
	scenario: string;
	before: Result;
	after: Result;
}

/**
 * v1 against v2 on the same harness, one call per op. v1 is given all three flags, its output closest to v2's String.
 * Every call site has been used before, as in a running program.
 */
export const comparison: Record<Group, Comparison[]> = {
	location: [
		{ scenario: 'Here, as a string', before: { ns: 436.5, allocs: 9 }, after: { ns: 72.7, allocs: 0 } },
		{ scenario: 'Caller(1), as a string', before: { ns: 449.5, allocs: 9 }, after: { ns: 93.5, allocs: 0 } },
		{ scenario: 'Here, into a buffer', before: { ns: 588.3, allocs: 9 }, after: { ns: 76, allocs: 0 } },
		{
			scenario: 'Here, as a string, 16 callers at once',
			before: { ns: 294.5, allocs: 9 },
			after: { ns: 6.3, allocs: 0 },
		},
	],
	errors: [{ scenario: 'fmt.Errorf with Here', before: { ns: 541.9, allocs: 12 }, after: { ns: 165.7, allocs: 3 } }],
	logging: [
		{ scenario: 'slog text line', before: { ns: 1680, allocs: 10 }, after: { ns: 562.3, allocs: 1 } },
		{ scenario: 'slog JSON line', before: { ns: 1110, allocs: 10 }, after: { ns: 430.5, allocs: 1 } },
	],
};

export interface SuiteResult {
	benchmark: string;
	what: string;
	ns: number;
	bytes: number;
	allocs: number;
}

/** v2's own benchmark suite, `go test -bench .`. */
export const suite: SuiteResult[] = [
	{ benchmark: 'Here', what: 'Capture the line that calls it', ns: 66.5, bytes: 0, allocs: 0 },
	{ benchmark: 'Caller', what: 'Capture the line one call up', ns: 87.1, bytes: 0, allocs: 0 },
	{ benchmark: 'HereString', what: 'Capture, then String', ns: 73.6, bytes: 0, allocs: 0 },
	{ benchmark: 'HereAppendText', what: 'Capture, then AppendText into a buffer', ns: 77.6, bytes: 0, allocs: 0 },
	{ benchmark: 'HereStringParallel', what: 'HereString from 16 goroutines at once', ns: 6.2, bytes: 0, allocs: 0 },
	{ benchmark: 'String', what: 'String of a frame already captured', ns: 7.1, bytes: 0, allocs: 0 },
	{ benchmark: 'AppendText', what: 'AppendText of a frame already captured', ns: 7.4, bytes: 0, allocs: 0 },
	{ benchmark: 'Errorf', what: 'fmt.Errorf("%s: %w", trace.Here(), err)', ns: 164.7, bytes: 104, allocs: 3 },
	{ benchmark: 'Slog/text', what: 'A slog text line with a frame attribute', ns: 598, bytes: 8, allocs: 1 },
	{ benchmark: 'Slog/JSON', what: 'A slog JSON line with a frame attribute', ns: 439.6, bytes: 8, allocs: 1 },
	{ benchmark: 'FirstUse', what: 'A call site\'s first String: resolved, kept, and the cache emptied again', ns: 190, bytes: 576, allocs: 7 },
];
