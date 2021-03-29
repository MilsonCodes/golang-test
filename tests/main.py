import threading, sys

target = 1000

#  Send the sequence 2, 3, 4, ... to list 'ch'.
def generate(ch):
    i = 1
    # for i in range(7920, 1, -1): # Must hack range in order to fill fast enough and still be able to use append
    while True:
        i += 1
        ch.insert(0, i)

# Copy the values from channel 'in' to channel 'out',
# removing those divisible by 'prime'.
def filter(inP, out, prime):
    while len(inP) > 0:
        i = inP.pop() # Receive value from 'in'.
        if i%prime != 0:
            out.insert(0, i) # Send 'i' to 'out'.

#  The prime sieve: Daisy-chain Filter processes.
if __name__=='__main__':
    ch = [] # Create a new channel.
    x = threading.Thread(target=generate, args=(ch, ), daemon=True)
    x.start()
    for i in range(target):
        # while (len(ch) < 1):
        #     wait()
        prime = ch.pop()
        print(prime)
        ch1 = []
        y = threading.Thread(target=filter, args=(ch, ch1, prime, ))
        y.start()
        ch = ch1.copy()
    print("1000th prime number: ", prime)
    x._delete()
    sys.exit()

