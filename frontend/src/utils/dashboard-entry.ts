// One-shot signal for the header entrance after authentication.
let pendingLoginEntry = false

export function markLoginEntry() {
  pendingLoginEntry = true
}

export function consumeLoginEntry() {
  const pending = pendingLoginEntry
  pendingLoginEntry = false
  return pending
}
