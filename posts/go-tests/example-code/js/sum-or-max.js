export function sumOrMax(a, b, max) {
  const sum = a + b;

  return sum < max ? sum : max;
}
