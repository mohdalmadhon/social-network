import assert from 'node:assert/strict'
import test from 'node:test'
import { normalizePostAudience } from './postAudience.js'

test('post audience uses account IDs, not array positions', () => {
  const followers = normalizePostAudience([
    { ID: 42, FirstName: 'Mira', LastName: 'North' },
    { ID: 87, FirstName: 'Theo', LastName: 'Vale' },
  ])

  assert.deepEqual(followers, [
    { id: 42, name: 'Mira North' },
    { id: 87, name: 'Theo Vale' },
  ])
})

test('post audience can also read an ID-keyed followers object', () => {
  const followers = normalizePostAudience({
    106: { FirstName: 'Aya', LastName: 'Benali' },
  })

  assert.deepEqual(followers, [{ id: 106, name: 'Aya Benali' }])
})

test('post audience ignores entries that do not contain a valid account ID', () => {
  assert.deepEqual(normalizePostAudience([{ FirstName: 'No ID' }]), [])
})
