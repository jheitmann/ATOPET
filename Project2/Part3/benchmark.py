from statistics import stdev, mean
from your_code import *
# First metric, size of the signature
def benchmark_size_of_signature():
    """
        This test aims at verifying that if the sigmas aren't conform they will be refused
    """
    server = Server()
    client = Client()

    #Generating the keys
    pk_serialized, sk_serialized = server.generate_ca("a,b,c")
    
    #Registering the user on the server

    issuance_request, private_state = client.prepare_registration(pk_serialized, "weewoo", "a,b,c")

    response = server.register(sk_serialized, issuance_request, "weewoo", "a,b,c")

    credential = client.proceed_registration_response(sk_serialized,response, private_state)
    counts = []

    stdevs = []
    means = []
    m = ""
    for i in range(0,1001,500):
        print(i)
        data = []
        #Trying every length 1000 times
        for k in range(1000):
            m = ("a"*i).encode()
            c = int.from_bytes(sha256(m).digest(), "big") % G1.order()

            #Trying to sign a message
            sig = client.sign_request(pk_serialized, credential, m,"a,b")
            data.append(len(sig))
            #Verifying the signature
            assert server.check_request_signature(pk_serialized, m, "a,b", sig) == True
        counts.append(data)
        means.append(mean(data))
        stdevs.append(stdev(data))
    print("--- Results for the size of the signature ---")
    print(means)
    print(stdevs)

benchmark_size_of_signature()