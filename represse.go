import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func writeWithTransactionBlock(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		row, err := txn.ReadRow(ctx, "Singers", spanner.Key{"Id": 1}, []string{"Name"})
		if err != nil {
			return err
		}
		var name string
		if err := row.ColumnByName("Name", &name); err != nil {
			return err
		}
		fmt.Fprintf(w, "%s\n", name)
		m := spanner.Update("Singers", []string{"Name", "SingerId"}, []interface{}{"Elena", 1})
		_, err = txn.Update(ctx, m)
		return err
	})
	return err
}
  
