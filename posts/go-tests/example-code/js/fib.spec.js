import { fib } from './fib';

describe('fib', () => {
  it.each([
    [0, 0],
    [1, 1],
    [2, 1],
    [3, 2],
    [4, 3],
    [5, 5],
    [6, 8],
    [7, 13],
    [8, 21],
  ])('fib(%i) should be %i', (n, want) => {
    const have = fib(n);

    expect(have).toBe(want);
  });
});
