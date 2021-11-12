resource "aiven_project" "test_kafka_project" {
	project = "my-kafka-project"
}

resource "aiven_kafka" "my-test-kafka" {
	project = aiven_project.test_kafka_project.project
	cloud_name = "google-europe-west1"
	plan = "startup-2"
	service_name = "my-test-kafka"
	kafka_user_config {
		kafka_version = "2.8"
		kafka_rest = true
	}
}

resource "aiven_kafka_topic" "topic-kiln" {
	topic_name = "kiln"
	partitions = 3
	service_name = aiven_kafka.my-test-kafka.service_name
	replication = 2
	project = aiven_project.test_kafka_project.project
}
