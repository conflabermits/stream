## Summary

Quick note on how to create a poll.

## Command

```text
$ curl -X POST 'https://api.twitch.tv/helix/polls' -H 'Authorization: Bearer sdfyo32rsdf89uewfrf23ejbjkwdys' -H 'Client-Id: d9023rnjklsdncvsdo9iwojrnnsdff' -H 'Content-Type: application/json' -d '{
  "broadcaster_id":"123456789",
  "title":"Heads or Tails, but cats?",
  "choices":[{
    "title":"lion"
  },
  {
    "title":"tiger"
  }],
  "channel_points_voting_enabled":true,
  "channel_points_per_vote":1000,
  "duration":60
}'
```

## Response

```text
{"data":[{"id":"12345678-abcd-6543-wxyz-098765432199","broadcaster_id":"123456789","broadcaster_name":"conflabermits","broadcaster_login":"conflabermits","title":"Heads or Tails, but cats?","choices":[{"id":"98765432-xxxx-1111-zzzz-123456667899","title":"lion","votes":0,"channel_points_votes":0},{"id":"12121212-3434-5656-7878-9a8b7c6d5e4f","title":"tiger","votes":0,"channel_points_votes":0}],"channel_points_voting_enabled":true,"channel_points_per_vote":1000,"status":"ACTIVE","duration":60,"started_at":"2024-06-27T03:02:13.712731778Z"}]}
```

## Browser Source URL

https://www.twitch.tv/popout/conflabermits/poll

## References

* https://dev.twitch.tv/docs/api/reference/#create-poll

