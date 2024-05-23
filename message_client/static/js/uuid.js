/*
 * MIT License
 *
 * Copyright (c) 2023 Fabio Lima
 * 
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 * 
 *  The above copyright notice and this permission notice shall be included in
 *  all copies or substantial portions of the Software.
 * 
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 *  THE SOFTWARE.
 */

// REFERENCES:
// https://datatracker.ietf.org/doc/html/rfc4122
// https://www.ietf.org/archive/id/draft-ietf-uuidrev-rfc4122bis-12.html

function bigRandom(bits) {
  // IEEE-754 mantissa: 52 bits
  if (bits > 52) { bits = 52 };
  // `Math.random()` is not cryptographically secure
  return BigInt(Math.floor(Math.random() * 2 ** bits));
}
    
function toUUIDString(bignum) {
    const digits = bignum.toString(16).padStart(32, "0");
    return `${digits.substring(0, 8)
          }-${digits.substring(8, 12)
          }-${digits.substring(12, 16)
          }-${digits.substring(16, 20)
          }-${digits.substring(20, 32)}`;
}

// UUIDv4 (Random)
function uuid4() {
  return toUUIDString(
    bigRandom(48) << 80n | // Random A
    (0x4n << 76n)        | // Version
    bigRandom(12) << 64n | // Random B
    (0b1n << 63n)        | // Variant
    bigRandom(14) << 48n | // Random C
    bigRandom(48)          // Random C
  );
}

// UUIDv7 (Time+Random)
function uuid7() {
  let milli = Date.now()
  return toUUIDString(
    BigInt(milli) << 80n | // Time
    (0x7n << 76n)        | // Version
    bigRandom(12) << 64n | // Random A
    (0b1n << 63n)        | // Variant
    bigRandom(14) << 48n | // Random B
    bigRandom(48)          // Random B
  );
}

// UUID_nil
function uuidNil() {
  return toUUIDString(
      BigInt(0) << 80n |  // Random A
      (0x0n << 76n)    |  // Version
      BigInt(0) << 64n |  // Random B
      (0b0n << 63n)    |  // Variant
      BigInt(0) << 48n |  // Random C
      BigInt(0)           // Random C
  );
}

// TESTS
/*
console.log("\n[UUIDv4]\n")
for (i = 0; i < 10; i++) {
    console.log(uuid4());
}

console.log("\n[UUIDv7]\n")
for (i = 0; i < 10; i++) {
    console.log(uuid7());
}

console.log("\n[UUID_nil]\n")
for (i = 0; i < 10; i++) {
    console.log(uuidNil());
}
*/