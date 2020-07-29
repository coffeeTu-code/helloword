import os
import sys
import io
import json
import argparse
import numpy as np
from scipy import misc
from multiprocessing import Pool
import imageio

def process(asset_dir, sample, lowpass=13):

    src_png = "{}/source/{}.png".format(asset_dir, sample['id'])

    w = int(sample['meta']['bandwidth'])
    h = int(sample['meta']['features'])

    counter = np.zeros((4,), dtype=np.int32)
    if not os.path.exists(src_png):

        src = np.zeros( (h, w), dtype=np.uint8)
        y = 0

        counter[1] += h*w

        for sdr in sample['sdr']:

            x_ = sdr[0]
            z_ = sdr[1]

            n = len(sdr[0]) if len(sdr[0]) < lowpass else lowpass

            for i in range(n):
                x = int(x_[i])
                z = int(z_[i])
                src[y, x] = z
                counter[2] += 1
                counter[3] += z

            y = y+1
        
        imageio.imwrite(src_png, src)
    else:
        counter[1] = 1
        counter[2] = -1
        counter[3] = 1

    tgt_png = "{}/target/{}.png".format(asset_dir, sample['id'])

    if not os.path.exists(tgt_png):

        tgt = np.zeros( (1,2,1), dtype=np.uint8)

        tgt[0,0,0] = 255*(1 - sample['target'])
        tgt[0,1,0] = 255*sample['target']

        imageio.imwrite(tgt_png, tgt)

    if(sample['target'] == 1):
        counter[0] += 1
        #print("Positive Sample ", sample['id'], " feature::", sample['digest'])
    
    return counter

parser = argparse.ArgumentParser(description='')
parser.add_argument('--source', dest='source', default="sample.jss", help='path for data source')
parser.add_argument('--output', dest='output', default="asset", help='path for data source')
parser.add_argument('--cutoff', dest='cutoff', type=int, default=-1, help='# of observations')
parser.add_argument('--lowpass', dest='lowpass', type=int, default=13, help='# of SDR digits per feature')
args = parser.parse_args()

def make_dir(path):
  if not os.path.exists(path):
    os.makedirs(path)

def line_process(line):
    encoded = json.loads(line)
    asset_dir = args.output
    counter = process(asset_dir, encoded, lowpass=args.lowpass)
    return counter

def main():

    filepath = args.source
    cutoff = args.cutoff

    if not os.path.isfile(filepath):
       print("File path {} does not exist. Exiting...".format(filepath))
       sys.exit()

    asset_dir = args.output

    make_dir("{}/source".format(asset_dir))
    make_dir("{}/target".format(asset_dir))
  
    with open(filepath) as fp:

        cnt = 0
        sums = np.zeros((4,), dtype=np.int64)

        with Pool(10) as pool:
            for x in pool.imap(line_process, (line.strip() for line in fp), 10):
                sums = sums + x
                cnt += 1
                if(cnt > 0 and cnt % 100 == 0):
                    print(cnt, " processed, positive ratio={}, sparsity={}, mean={}".format(sums[0]/cnt, sums[2]/sums[1], sums[3]/sums[2]))

        print("________DONE__________")
        print(cnt, " processed, positive ratio={}, sparsity={}, mean={}".format(sums[0]/cnt, sums[2]/sums[1], sums[3]/sums[2]))

if __name__ == '__main__':
   main()