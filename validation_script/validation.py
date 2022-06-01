import subprocess
from multiprocessing.pool import ThreadPool
import os
from datetime import datetime


command = './validation_script -timing -benchmark={benchmark} -magic-memory-copy'

benchmarks = [
    'aes',
    'atax',
    'bicg',
    'bitonicsort',
    'conv2d',
    'fastwalshtransform',
    'fir',
    'fft',
    'floydwarshall',
    'im2col',
    'kmeans',
    'matrixmultiplication',
    'matrixtranspose',
    'nbody',
    'nw',
    'pagerank',
    'relu',
    'simpleconvolution',
    'spmv',
    'stencil2d',
]

results_dir = ''


def run_exp(benchmark, projection):
    try:
        stdout_file = benchmark
        stderr_file = benchmark

        if projection:
            stdout_file += '_projection'
            stderr_file += '_projection'

        stdout_file = f'{results_dir}{stdout_file}.stdout'
        stderr_file = f'{results_dir}{stderr_file}.stderr'

        cmd = command.format(benchmark=benchmark)
        if projection:
            cmd += ' -allow-projection'

        with open(stdout_file, 'w') as stdout, open(stderr_file, 'w') as stderr:
            process = subprocess.Popen(
                cmd,
                shell=True,
                stdout=stdout, stderr=stderr,
            )
            process.wait()
    except Exception as e:
        print(e)


def create_results_dir():
    global results_dir

    dir_name = datetime.now().strftime("%Y-%m-%d-%H-%M-%S")
    results_dir = 'results/' + dir_name + '/'

    if not os.path.exists('results'):
        os.makedirs('results')

    if not os.path.exists(results_dir):
        os.makedirs(results_dir)


def main():
    create_results_dir()

    process = subprocess.Popen("go build", shell=True)
    process.wait()

    tp = ThreadPool()
    for exp in benchmarks:
        tp.apply_async(run_exp, args=(exp, False,))
        tp.apply_async(run_exp, args=(exp, True,))

    tp.close()
    tp.join()


if __name__ == '__main__':
    main()
