# Cloud Foundry Cassandra Service Broker

A BOSH release of a Cassandra Service Broker for Cloud Foundry (does not contain cassandra itself).

## Register the Service Broker<a name="register_broker"></a>

### BOSH errand

BOSH errands were introduced in version 2366 of the BOSH CLI, BOSH Director, and stemcells.

```
$ bosh run errand broker-registrar
```

Note: the broker-registrar errand will fail if the broker has already been registered, and the broker name does not match the manifest property `jobs.broker-registrar.properties.broker.name`. Use the `cf rename-service-broker` CLI command to change the broker name to match the manifest property then this errand will succeed.

### Manually

1. First register the broker using the `cf` CLI.  You must be logged in as an admin.

    ```
    $ cf create-service-broker apache-cassandra BROKER_USERNAME BROKER_PASSWORD URL
    ```

    `BROKER_USERNAME` and `BROKER_PASSWORD` are the credentials Cloud Foundry will use to authenticate when making API calls to the service broker. Use the values for manifest properties `jobs.cf-cassandra-broker.properties.auth_username` and `jobs.cf-cassandra-broker.properties.auth_password`.

    `URL` specifies where the Cloud Controller will access the Cassnadra broker. Use the value of the manifest property `jobs.cf-cassnadra-broker.properties.external_host`.

    For more information, see [Managing Service Brokers](http://docs.cloudfoundry.org/services/managing-service-brokers.html).

2. Then [make the service plan public](http://docs.cloudfoundry.org/services/managing-service-brokers.html#make-plans-public).


## De-register the Service Broker<a name="deregister_broker"></a>

The following commands are destructive and are intended to be run in conjuction with deleting your BOSH deployment.

### BOSH errand

BOSH errands were introduced in version 2366 of the BOSH CLI, BOSH Director, and stemcells.

This errand runs the two commands listed in the manual section below from a BOSH-deployed VM. This errand should be run before deleting your BOSH deployment. If you have already deleted your deployment follow the manual instructions below.

```
$ bosh run errand broker-deregistrar
```

### Manually

Run the following:

```
$ cf purge-service-offering apache-cassandra
$ cf delete-service-broker apache-cassandra
```

## Security Groups<a name="register_broker"></a>

Since [cf-release](https://github.com/cloudfoundry/cf-release) v175, applications by default cannot to connect to IP addresses on the private network. This prevents applications from connecting to the Cassnadra service. To enable access to the service, create a new security group for the IPs configured in your manifest for the Cassnadra cluster.

1. Add the rule to a file in the following json format; multiple rules are supported.

  ```
  [
      {
        "destination": "192.168.111.30-192.168.111.34",
        "protocol": "tcp",
        "ports": "9042"
      }
  ]
  ```
- Create a security group from the rule file.
  <pre class="terminal">
  $ cf create-security-group cassandra rule.json
  </pre>
- Enable the rule for all apps
  <pre class="terminal">
  $ cf bind-running-security-group cassandra
  </pre>

Changes are only applied to new application containers; in order for an existing app to receive security group changes it must be restarted.
