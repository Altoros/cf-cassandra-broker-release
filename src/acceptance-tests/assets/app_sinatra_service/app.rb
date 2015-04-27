require "sinatra"
require "cassandra"
require "json"
require "securerandom"

get "/test" do
  services = JSON.parse(ENV["VCAP_SERVICES"])
  credentials = services["apache-cassandra"][0]["credentials"]
  begin
    cluster = Cassandra.cluster({
      username: credentials["username"],
      password: credentials["password"],
      hosts:    credentials["nodes"],
      port:     credentials["cql_port"]
    })
    session = cluster.connect(credentials["keyspace"])
    session.execute("SELECT now() FROM system.local")
    "works"
  rescue => e
    e.inspect
  end
end
