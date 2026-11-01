---
name: Object storage (buckets)
keywords: [bucket, buckets, object storage, s3, files, file upload, upload, images, pictures, photos, media, documents, pdf]
status: draft
---
- One bucket per service, created by the platform. The name comes from the environment
  variable `STORAGE_BUCKET`. Never hardcode bucket names.
- Use the S3-compatible API via the AWS SDK for Go v2; credentials come from the
  workload identity, never from keys in code or config.
- Object keys: `<service>/<entity>/<id>/<filename>`, for example
  `catvote/pictures/42/cat.jpg`.
- Serve objects through your service or a pre-signed URL that expires within 15 minutes.
  Buckets are never public.

Copying images from the internet into a bucket: first read `goldenpath lookup external-services`.
