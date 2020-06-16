resource "aws_s3_bucket" "codepipeline_bucket" {
  bucket = "flipbot-ci-bucket"
  acl    = "private"
}

resource "aws_s3_bucket" "codepipeline_bucket_log" {
  bucket = "flipbot-ci-bucket-log"
  acl    = "private"
}


resource "aws_iam_role" "codepipeline_role" {
  name = "build-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "codepipeline.amazonaws.com",
        "Service": "codebuild.amazonaws.com"
      },
      "Action": [
          "sts:AssumeRole"
        ]
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "codepipeline_policy" {
  name = "codepipeline_policy"
  role = aws_iam_role.codepipeline_role.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect":"Allow",
      "Action": [
        "s3:GetObject",
        "s3:GetObjectVersion",
        "s3:GetBucketVersioning",
        "s3:PutObject"
      ],
      "Resource": [
        "${aws_s3_bucket.codepipeline_bucket.arn}",
        "${aws_s3_bucket.codepipeline_bucket.arn}/*",
        "${aws_s3_bucket.codepipeline_bucket_log.arn}",
        "${aws_s3_bucket.codepipeline_bucket_log.arn}/*"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "codebuild:BatchGetBuilds",
        "codebuild:StartBuild",
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_codepipeline" "codepipeline" {
  name     = "flipbot-pipeline"
  role_arn = aws_iam_role.codepipeline_role.arn

  artifact_store {
    location = aws_s3_bucket.codepipeline_bucket.bucket
    type     = "S3"
  }

  stage {
    name = "Source"

    action {
      name             = "Source"
      category         = "Source"
      owner            = "ThirdParty"
      provider         = "GitHub"
      version          = "1"
      output_artifacts = ["source_output"]

      configuration = {
        Owner  = "ugjka"
        Repo   = "flipbot"
        Branch = "master"
      }
    }
  }
  stage {
    name = "Build"

    action {
      name             = "Build"
      category         = "Build"
      owner            = "AWS"
      provider         = "CodeBuild"
      input_artifacts  = ["source_output"]
      output_artifacts = ["build_output"]
      version          = "1"

      configuration = {
        ProjectName = aws_codebuild_project.flipbot.name
      }
    }
  }
}

resource "aws_codebuild_project" "flipbot" {
  name          = "build"
  description   = "build flipbot"
  build_timeout = "10"
  service_role  = aws_iam_role.codepipeline_role.arn

  artifacts {
    type = "CODEPIPELINE"
  }

  environment {
    compute_type                = "BUILD_GENERAL1_SMALL"
    image                       = "_/golang:1.13.12"
    type                        = "LINUX_CONTAINER"
    image_pull_credentials_type = "CODEBUILD"
  }

  logs_config {
    s3_logs {
      status   = "ENABLED"
      location = "${aws_s3_bucket.codepipeline_bucket_log.id}/build-log"
    }
  }

  source {
    type = "CODEPIPELINE"
  }
}
