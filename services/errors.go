// Most functions shouldn't return errors wrapped with the package name.
// It only really makes sense if a function introduces its own errors but isn't a Start or Shutdown method which log.Fatalf or log a warning instead.

package services
