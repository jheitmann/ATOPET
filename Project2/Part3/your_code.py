"""
Classes that you need to complete.
"""

# Optional import
from serialization import jsonpickle

# Imports for PS scheme
from petrelic.multiplicative.pairing import G1,G2, GT 

# Imports for NIZK
from binascii import hexlify
from hashlib import sha256

# Constants
L = 4

class Server:
    """Server"""

    @staticmethod
    def generate_ca(valid_attributes):
        """Initializes the credential system. Runs exactly once in the
        beginning. Decides on schemes public parameters and choses a secret key
        for the server.

        Args:
            valid_attributes (string): a list of all valid attributes. Users cannot
            get a credential with a attribute which is not included here.

            Note: You can use JSON to encode valid_attributes in the string.

        Returns:
            (tuple): tuple containing:
                byte[] : server's pubic information
                byte[] : server's secret key
            You are free to design this as you see fit, but all commuincations
            needs to be encoded as byte arrays.
        """

        # Code to generate the secret key
        private_key_array = []
        public_key_array = []
        ys_in_G1 = []
        g = G2.generator()
        for i in range(L+1):
            c = G1.order().random()
            private_key_array.append(c)
            public_key_array.append(g ** c)
            ys_in_G1.append(G1.generator() ** c)
        x = private_key_array[0]
        X = public_key_array[0]
        X_G1 = G1.generator() ** x
        y_s = private_key_array[1:]
        Y_s = public_key_array[1:]
        ys_in_G1 = ys_in_G1[1:]

        # Some sanity checks on the key
        for i in range(len(y_s)):
            if not (G1.generator().pair(Y_s[i]) == ys_in_G1[i].pair(G2.generator())):
                raise Exception("Key generation failed !")


        public_information = {"Ys" : Y_s, "X_G1" : X_G1 ,"X" : X,"Ys_G1":ys_in_G1, "validAttributes" : valid_attributes.split(",")}
        pk = jsonpickle.encode(public_information).encode()
        private_key =        {"Ys" : Y_s, "X_G1" : X_G1 ,"X" : X,"Ys_G1":ys_in_G1, "validAttributes" : valid_attributes.split(","),"x": x,  "ys": y_s}
        return pk, jsonpickle.encode(private_key).encode()

    def register(self, server_sk, issuance_request, username, attributes):
        """ Registers a new account on the server.

        Args:
            server_sk (byte []): the server's secret key (serialized)
            issuance_request (bytes[]): The issuance request (serialized)
            username (string): username
            attributes (string): attributes

            Note: You can use JSON to encode attributes in the string.

        Return:
            response (bytes[]): the client should be able to build a credential
            with this response.
        """

        #Reconstructing the secret key
        sk = jsonpickle.decode(server_sk)
        #Decoding the request
        decoded_request = jsonpickle.decode(issuance_request)
        C = decoded_request["C"]

        #Verifying the value for C
        s_sk = decoded_request["s_sk"]
        s_t  = decoded_request["s_t"]
        R    = decoded_request["R"]
        c = int.from_bytes(sha256((username + attributes).encode()).digest(), "big") % G1.order() # hashing 
        
        if (sk["Ys_G1"][0]**s_sk) * (G1.generator() ** s_t) * C**c != R:
            raise Exception("Invalid register !")

        #Generating a random u
        u = G1.order().random()
        g_pow_u = G1.generator() ** u
        X = sk["X_G1"] #Note : I decided to have X in the secret key, rather than recompute it everytime

        #Building the Ys
        Ys = sk["Ys_G1"]

        sigma_2 = X*C
        #The subscription attributes
        for att in attributes.split(","):
            if att not in sk["validAttributes"]:
                raise ValueError("Non-recognized attribute")
            else:
                index_of_att = sk["validAttributes"].index(att)
                y_i = Ys[index_of_att + 1] #+1 because of the hidden attribute
                sigma_2 = sigma_2 * y_i
        json_res = {"sigma1" : g_pow_u, "sigma2": sigma_2 ** u}
        return jsonpickle.encode(json_res).encode() 

    def check_request_signature(
        self, server_pk, message, revealed_attributes, signature
    ):
        """

        Args:
            server_pk (byte[]): the server's public key (serialized)
            message (byte[]): The message to sign
            revealed_attributes (string): revealed attributes
            signature (bytes[]): user's autorization (serialized)

            Note: You can use JSON to encode revealed_attributes in the string.

        Returns:
            valid (boolean): is signature valid
        """
        sig = jsonpickle.decode(signature)
        att = revealed_attributes.split(",")
        pk = jsonpickle.decode(server_pk)
        valid_attributes = pk["validAttributes"]

        R = sig["R"]
        c = sig["c"]
        s_is = sig["s_is"]
        sigma1 = sig["sigma"][0]
        sigma2 = sig["sigma"][1]
        if sigma1 == G1.generator() ** G1.order(): #Checking that sigma1 != 1
            return False

        # Verifying the hash was legit
        c_us = int.from_bytes(sha256(message).digest(), "big") % G1.order()
        if(c_us != c):
            return False

        # Need to recreate the value the proof of knowledge is showing
        acc = sigma2.pair(G2.generator()) / sigma1.pair(pk["X"])
        for i in range(len(att)):
            acc = acc / (sigma1).pair(pk["Ys"][valid_attributes.index(att[i])+1])

        values_to_prove = [] # Storing the values used in the proof of knowledge
        values_to_prove.append((sigma1).pair(G2.generator()))
        values_to_prove.append(sigma1.pair(pk["Ys"][0]))
        for att in valid_attributes:
            if att not in revealed_attributes.split(","):
                values_to_prove.append(sigma1.pair(pk["Ys"][valid_attributes.index(att)+1]))
        thing_to_add_to_r = acc ** c_us
        for v,exp in zip(values_to_prove, s_is):
            thing_to_add_to_r = thing_to_add_to_r * (v**exp)
        return R == thing_to_add_to_r


class Client:
    """Client"""

    def prepare_registration(self, server_pk, username, attributes):
        """Prepare a request to register a new account on the server.

        Args:
            server_pk (byte[]): a server's public key (serialized)
            username (string): username
            attributes (string): user's attributes

            Note: You can use JSON to encode attributes in the string.

        Return:
            tuple:
                byte[]: an issuance request
                (private_state): You can use state to store and transfer information
                from prepare_registration to proceed_registration_response.
                You need to design the state yourself.
        """

        t = G1.order().random()
        pk = jsonpickle.decode(server_pk)
        Ys_in_G1 = pk["Ys_G1"]

        random_sk = G1.order().random() # Some secret key as asked
        C = (G1.generator() ** t) * (Ys_in_G1[0] ** random_sk)
        for att in attributes.split(","):
            if att not in pk["validAttributes"]:
                raise ValueError("Attribute not in list of accepted attribute by the server")


        issuance_request = {}
        issuance_request["C"] = C
        r_sk = G1.order().random()
        r_t = G1.order().random()
        c = int.from_bytes(sha256((username  + attributes).encode()).digest(), "big") % G1.order() # hashing 
        
        s_sk = ((r_sk - random_sk*c) + G1.order()) % G1.order()
        s_t  = ((r_t - t*c) + G1.order()) % G1.order()

        R = (Ys_in_G1[0] ** r_sk) * (G1.generator() ** r_t)
        if (Ys_in_G1[0] ** s_sk) * (G1.generator() ** s_t) * C**c != R:
            raise Exception("Invalid signature !")
        
        issuance_request["s_t"] = s_t
        issuance_request["s_sk"] = s_sk
        issuance_request["R"] = R


        return jsonpickle.encode(issuance_request).encode(), (random_sk,t) #t is the inner state, we only need to remember that (since the key as secret attribute is completely useless)
        
    def proceed_registration_response(self, server_pk, server_response, private_state):
        """Process the response from the server.

        Args:
            server_pk (byte[]): a server's public key (serialized)
            server_response (byte[]): the response from the server (serialized)
            private_state (private_state): state from the prepare_registration
            request corresponding to this response

        Return:
            credential (byte []): create an attribute-based credential for the user
        """
        
        sigma = jsonpickle.decode(server_response)
        (random_sk, t) = private_state
        sigma_bis = (sigma["sigma1"],sigma["sigma2"]/(sigma["sigma1"] ** t))

        #Verifying the received signature
        pk = jsonpickle.decode(server_pk)
        Y_s = pk["Ys"]
        acc = G2.neutral_element()
        for (i, y_i) in zip(range(len(Y_s)), Y_s):
            if i == 0:
                acc = acc * (y_i ** random_sk)
            else:
                acc = acc * y_i
        #Checking that the signature is valid
        if not (sigma_bis[0].pair(acc*pk["X"]) == sigma_bis[1].pair(G2.generator())):
            raise Exception("Invalid signature !")

        return jsonpickle.encode({"sigma": (sigma_bis[0], sigma_bis[1]), "random_sk":random_sk}).encode()


    def sign_request(self, server_pk, credential, message, revealed_info):
        """Signs the request with the clients credential.

        Arg:
            server_pk (byte[]): a server's public key (serialized)
            credential (byte[]): client's credential (serialized)
            message (byte[]): message to sign
            revealed_info (string): attributes which need to be authorized

            Note: You can use JSON to encode revealed_info.

        Returns:
            byte []: message's signature (serialized)
        """
        # Generating r,t
        r = G1.order().random()
        t = G1.order().random()
        
        #Decoding the public key
        pk = jsonpickle.decode(server_pk)
        valid_attributes = pk["validAttributes"]

        decoded_credential = jsonpickle.decode(credential)
        sigma = decoded_credential["sigma"]
        random_sk = decoded_credential["random_sk"]

        sigma_bis = (sigma[0] ** r,  ((sigma[0] ** t) * sigma[1]) ** r)

        values_to_prove = [] # Storing the values used in the proof of knowledge

        com = (sigma_bis[0]).pair(G2.generator()) ** t
        values_to_prove.append((sigma_bis[0]).pair(G2.generator()))
        # The secret part of the multiplication
        com = com * sigma_bis[0].pair(pk["Ys"][0] ** random_sk) 
        values_to_prove.append(sigma_bis[0].pair(pk["Ys"][0]))
        # Attributes
        for att in valid_attributes:
            if att not in revealed_info.split(","):
                #raise ValueError("Attribute not in the possible ones !")
                com = com * sigma_bis[0].pair(pk["Ys"][valid_attributes.index(att)+1])
                values_to_prove.append(sigma_bis[0].pair(pk["Ys"][valid_attributes.index(att)+1]))
        c = int.from_bytes(sha256(message).digest(), "big") % G1.order() # hashing 
        

        #Adding the signature
        r_is = []
        R = GT.generator() ** GT.order() # The neutral element of the group
        for value in values_to_prove:
            r_i = G1.order().random()
            r_is.append(r_i)
            R = R * (value ** r_i) # creating the R part of the protocol

        

        s_is = []
        for (i,r_i) in zip(range(len(r_is)),r_is):
            tmp = r_i - c % G1.order()
            if i == 0:
                tmp = (r_i - c*t) % G1.order()
            if i == 1:
                tmp = (r_i - c*random_sk) % G1.order()


            if tmp < 0:
                tmp += G1.order()
            s_is.append(tmp)


        
    

        proof = {"c": c, "R":R, "s_is":s_is, "sigma": sigma_bis}
        return jsonpickle.encode(proof).encode()
