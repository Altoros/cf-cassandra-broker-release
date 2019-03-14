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
        "ports": "9042,9160"
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

## Running Acceptance Tests

To run the Cassandra acceptance tests you will need:
- a running CF instance
- credentials for a CF Admin user
- a deployed Cassandra Broker Release and the plan made public
- a security group granting access to the service for applications

### Using BOSH errands

BOSH errands were introduced in version 2366 of the BOSH CLI, BOSH Director, and stemcells.

The following properties must be included in the manifest (most will be there by default):
- cf.api_url:
- cf.admin_username:
- cf.admin_password:
- cf.apps_domain:
- cf.skip_ssl_validation:
- broker.host:

```
$ bosh run errand acceptance-tests
```

### Manually

To run the acceptance tests manually you will also need an environment variable `$CONFIG` which points to a `.json` file that contains the application domain.

1. Install `go` by following the directions found [here](http://golang.org/doc/install)
2. `cd` into `cf-cassandra-broker-release/src/acceptance-tests/`
3. Update `cf-cassandra-broker-release/src/acceptance-tests/integration_config.json`

   The following commands provide a shortcut to configuring `integration_config.json` with values for a [bosh-lite](https://github.com/cloudfoundry/bosh-lite)
deployment. Copy and paste this into your terminal, then open the resulting `integration_config.json` in an editor to replace values as appropriate for your environment.

    ```bash
    cat > integration_config.json <<EOF
    {
      "api":                 "api.10.244.0.34.xip.io",
      "admin_user":          "admin",
      "admin_password":      "admin",
      "apps_domain":         "10.244.0.34.xip.io",
      "service_name":        "apache-cassandra",
      "plan_name":           "free",
      "broker_host":         "cassandra.10.244.0.34.xip.io",
      "skip_ssl_validation": true
    }
    EOF
    export CONFIG=$PWD/integration_config.json
    ```

    Note: `skip_ssl_validation` requires CLI v6.0.2 or newer.

4. Run  the tests

  ```
  $ ./bin/test
  ```

### NOTE: 
If somehow bosh will not follow symlinks (packages/common;src/common) while you creating bosh-release - then just replace the links with actual files.

