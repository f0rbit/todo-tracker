

// write a test case for adding 3 numbers together
function add(a: number, b: number, c: number): number {
  return a + b + c;
}

/**
 * @todo we need to import the `test` and `excect` functions
 */
test('add 3 numbers together', () => {
  expect(add(1, 2, 3)).toBe(6);
});

