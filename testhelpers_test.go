package allora

// testMnemonic is a well-known BIP-39 test vector shared across the package's tests
// (remote_signer_test.go, tx_test.go). It lives in this dedicated helper file so no test
// file silently depends on a constant declared in another, per the team Go style guide.
const testMnemonic = "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
