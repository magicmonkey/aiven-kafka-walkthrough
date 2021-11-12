/*
This generates an authenticated Kafka client.  Code is based on
example at https://help.aiven.io/en/articles/5344122-go-examples-for-testing-aiven-for-apache-kafka
*/

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/Shopify/sarama"
)

func getKafkaClient() sarama.Client {
	keypair, err := tls.LoadX509KeyPair(serviceCertPath, serviceKeyPath)
	if err != nil {
		fmt.Println("If you have an invalid cert/key, please get a new one using:")
		fmt.Println("  avn service user-creds-download --username avnadmin <kafka-service-name>")
		fmt.Println("")
		panic(err)
	}

	caCert, err := ioutil.ReadFile(projectCaPath)
	if err != nil {
		panic(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{keypair},
		RootCAs:      caCertPool,
	}

	// init config, enable errors and notifications
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Net.TLS.Enable = true
	config.Net.TLS.Config = tlsConfig
	config.Version = sarama.V0_10_2_0

	// init producer
	brokers := []string{serviceURI}

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		panic(err)
	}

	return client
}
