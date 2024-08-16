# github.com/cyrusaf/page

Turn paginated APIs into Golang iterators.

## Usage

All a `page.Iter[I, P]` requires is that you define a `page.Read[I, P]` function.
`page.Iter` will then invoke the read function whenever it hits the end of the given page.

```golang
type Read[I, P any] func(ctx context.Context, nextPage *P) ([]I, *P, error)
```

Then create a new `page.Iter` using your `read` function. For example:

```golang
for item, err := range page.Iter(ctx, readPage) {
    // ...
}

func readPage(ctx context.Context, nextPage *string) ([]int, *string, error) {
    resp, err := getItems(ctx, nextPage) // getItems might be a paginated API call
    if err != nil {
        return nil, nil, err
    }
    return resp.Items, resp.NextPage, nil
}
```



## Example using the AWS SDK

```golang
func run(ctx context.Context) error {
    // Set up AWS client
    cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
    if err != nil {
        return fmt.Errorf("unable to load SDK config, %w", err)
    }
    ivsrt := ivsrealtime.NewFromConfig(cfg)

    // Create a page.Iter[types.StageSummary, string].
    // Pass in page.Read[I, T] = readStagePage() defined below.
    for stage, err := range page.Iter(ctx, readStagePage(ivsrt)) {
        if err != nil {
            return err
        }
        fmt.Printf("Stage ARN: %+v\n", *stage.Arn)
    }
    return nil
}

// Define page.Read[I, P], where I is the iterator item type and P is the page type.
// In this case, I=types.StageSummary, P=string.
func readStagePage(ivsrt *ivsrealtime.Client) func(context.Context, *string) ([]types.StageSummary, *string, error) {
    return func(ctx context.Context, nextPage *string) ([]types.StageSummary, *string, error) {
        resp, err := ivsrt.ListStages(ctx, &ivsrealtime.ListStagesInput{
            MaxResults: aws.Int32(100),
            NextToken:  nextPage,
        })
        if err != nil {
            return nil, nil, fmt.Errorf("listing stages: %w", err)
        }
    return resp.Stages, resp.NextToken, nil
    }
}
```
