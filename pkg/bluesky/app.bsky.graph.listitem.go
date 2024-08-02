package bluesky

// {
//   "lexicon": 1,
//   "id": "app.bsky.graph.listitem",
//   "defs": {
//     "main": {
//       "type": "record",
//       "description": "Record representing an account's inclusion on a specific list. The AppView will ignore duplicate listitem records.",
//       "key": "tid",
//       "record": {
//         "type": "object",
//         "required": ["subject", "list", "createdAt"],
//         "properties": {
//           "subject": {
//             "type": "string",
//             "format": "did",
//             "description": "The account which is included on the list."
//           },
//           "list": {
//             "type": "string",
//             "format": "at-uri",
//             "description": "Reference (AT-URI) to the list record (app.bsky.graph.list)."
//           },
//           "createdAt": { "type": "string", "format": "datetime" }
//         }
//       }
//     }
//   }
// }

// FromBskyGraphListitem is
func FromBskyGraphListitem() {

}

// ToBskyGraphListitem is
func ToBskyGraphListitem() {

}
