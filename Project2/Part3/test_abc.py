from your_code import *
def test_signature_works():
    """
        This test verifies that a normal user will have functionnality
    """
    server = Server()
    client = Client()

    #Generating the keys
    pk_serialized, sk_serialized = server.generate_ca("a,b,c")
    
    #Registering the user on the server
    issuance_request, private_state = client.prepare_registration(pk_serialized, "weewoo", "a,b,c")

    response = server.register(pk_serialized, issuance_request, "weewoo", "a,b,c")

    credential = client.proceed_registration_response(sk_serialized,response, private_state)

    #Trying to sign a message
    m = b"some message for test"
    sig = client.sign_request(pk_serialized, credential, m,"a,b")
    
    #Verifying the signature
    assert server.check_request_signature(pk_serialized, m, "a,b", sig) == True


from petrelic.multiplicative.pairing import G1,G2, GT 
def test_error_condition_respected():
    """
        This test aims at verifying that if the sigmas aren't conform they will be refused
    """
    server = Server()
    client = Client()

    #Generating the keys
    pk_serialized, sk_serialized = server.generate_ca("a,b,c")
    
    #Registering the user on the server

    m = b"some message for test"
    c = int.from_bytes(sha256(m).digest(), "big") % G1.order()

    credential = jsonpickle.encode({"R":3, "c":c, "sigma": (G1.generator(), G1.generator()), "random_sk": 1})
    #Trying to sign a message
    sig = client.sign_request(pk_serialized, credential, m,"a,b")
    
    #Verifying the signature
    assert server.check_request_signature(pk_serialized, m, "a,b", sig) == False

def test_correct_pk():
    """
        This test verifies that a normal user will not be fooled by the wrong server
    """
    server = Server()
    client = Client()

    #Generating the keys
    pk_serialized1, sk_serialized1 = server.generate_ca("a,b,c")
    pk_serialized2, sk_serialized2 = server.generate_ca("a,b,c")
    
    #Registering the user on the server
    issuance_request, private_state = client.prepare_registration(pk_serialized1, "weewoo", "a,b,c")

    try:
        response = server.register(sk_serialized2, issuance_request, "weewoo", "a,b,c")
        raise Exception("Should have otherwise ...")
    except Exception as e:
        assert str(e) == ("Invalid register !")


def test_random_sigma():
    """
        This test aims at verifying that if the sigmas aren't conform they will be refused
    """
    server = Server()
    client = Client()

    #Generating the keys
    pk_serialized, sk_serialized = server.generate_ca("a,b,c")
    
    #Registering the user on the server

    m = b"some message for test"
    c = int.from_bytes(sha256(m).digest(), "big") % G1.order()

    credential = jsonpickle.encode({"R":3, "c":c, "sigma": (G1.generator() ** G1.order().random(), G1.generator() ** G1.order().random()), "random_sk": 1})
    #Trying to sign a message
    sig = client.sign_request(pk_serialized, credential, m,"a,b")
    
    #Verifying the signature
    assert server.check_request_signature(pk_serialized, m, "a,b", sig) == False