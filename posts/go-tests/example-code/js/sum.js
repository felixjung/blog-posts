function isString(value) {
  return typeof value === 'string';
}

export function sum(a, b) {
  if (isString(a) || isString(b)) {
    return null;
  }

  return a + b;
}
