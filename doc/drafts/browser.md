# Browser

For the purposes of creating a bare, minimal browser of the event data that is uploaded by the
Hoover, a specification of the features of the app must be created.

## General Considerations

For reasons of simplicity, the browser will render detail of content for posts only, and other
event types will just be a simple hider that can be clicked to reveal the raw tag data of the
event.

The primary view will be controlled by specifying a range of timestamps, by picking a "newest"
timestamp and a value of a number of seconds of events to display.

The browser's main controls will be a set of 3 buttons that enable and disable the display of
matching events of the Nostr, Bluesky and Farcaster protocols, that will trigger a reload of the
page.

The Forward and Back buttons will step forward and backwards so that if the span is set to 10
seconds, each "forward" will increase the timestamp by 10 seconds and cover the next, exclusive
10 seconds of events.

Only the "post" events will render the post, with a visible user identity, and clicking on this
identity button will switch to a filter that shows only this user's events. All other event
types also the user identity will be clickable to show only their events (plus the time window
in force).

All events other than post will have a hider that when clicked reveals the raw event tag fields,
with the poster and timestamp as the leaders on the line.

All events will have a cryptographic verification that ensures that the public key of the poster
and the event ID hash match up with the event signature. This will produce a checkmark icon that
appears at the end of the top row for each event.

## Rationale

The primary purpose is to create a simple baseline example that shows the key elements of
building a basic viewer for Hoover events, without becoming excessively complex and costing more
time in development, and fulfilling the requirements as specified in
the [milestone document](../milestones.md).
