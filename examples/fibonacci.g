// Recursive Fibonnaci implementation.
print "Printing the 10-th Fibonnaci number";

fun fib(n) {
  if (n <= 1) {
    return n;
  }
  return fib(n-1) + fib(n-2);
}

print fib(10);