import { sumOrMax } from './sum-or-max';

describe('sumOrMax', () => {
  describe('when the sum is below the maximum', () => {
    it('should return the sum', () => {
      const a = 1;
      const b = 2;
      const max = 4;

      const want = 3;

      const have = sumOrMax(a, b, max);

      expect(have).toBe(want);
    });
  });

  describe('when the sum is at or above the maximum', () => {
    it('should return the maximum', () => {
      const a = 5;
      const b = 10;
      const max = 4;

      const want = max;

      const have = sumOrMax(a, b, max);

      expect(have).toBe(want);
    });
  });
});
