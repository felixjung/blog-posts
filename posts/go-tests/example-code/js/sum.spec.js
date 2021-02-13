import { sum } from './sum';

test('sum', () => {
  const a = 1;
  const b = 2;

  const want = 3;

  const got = sum(a, b);

  expect(got).toBe(want);
});
