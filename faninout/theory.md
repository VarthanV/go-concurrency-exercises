# Fan-out, Fan-In

- Sometimes stages in a pipeline can be particularly computationally expensive , When this happens upstream stages in the pipeline can become blocked while waiting for the expesive stages to complete.

- Not only that but the pipeline itself can take a long time to execute on whole.

- One of the interesting properties of pipeline is the ability they give to operate on the stream of data using`` combination of seperate, reorderable stages``. Can even reuse stages of pipeline multiple times.

- It would be interesting to reuse a single stage of the pipeline multiple times in an attempt to parallelize pulls from upstream stage.

**Fan-Out**: It is a term to describe the process of starting multiple goroutines to handle input from pipeline.

**Fan-In**: Process of combining multiple results into single channel

- This pattern is suited if the below rules are applicable
    - It doesn't rely on values that the stage had calculated before.
    - It takes a long time to run.

