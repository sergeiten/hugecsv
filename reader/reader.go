package reader

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	pb "github.com/sergeiten/hugecsv"
	"google.golang.org/grpc"
)

type Reader struct {
	filename string
	conn     *grpc.ClientConn
	client   pb.HugeCSVClient
	f        *os.File
}

func New(filename string) *Reader {
	return &Reader{
		filename: filename,
	}
}

func (r *Reader) Serve(ctx context.Context) error {
	err := r.init()
	if err != nil {
		return err
	}

	stream, err := r.client.Send(context.Background())
	if err != nil {
		return err
	}

	r.f, err = os.Open(r.filename)
	if err != nil {
		return err
	}

	br := bufio.NewReader(r.f)

	// ignore first line
	_, _ = br.ReadString('\n')

	done := make(chan bool)

	go func() {
		for {
			data, err := br.ReadString('\n')
			if err == io.EOF {
				done <- true
				return
			}

			if err != nil {
				pb.LogPrint(err, "failed to read string")
				continue
			}

			item := r.stringToItem(data)

			//fmt.Printf("%+v\n", item)

			err = stream.Send(item)
			if err != nil {
				log.Printf("%v.Send(%v) = %v\n", stream, item, err)
			}

			select {
			case <-ctx.Done():
				fmt.Println("Received shut down signal")
				done <- true
				return
			default:
			}
		}
	}()

	<-done

	defer r.f.Close()
	defer r.conn.Close()

	summary, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}

	log.Printf("Send summary: %v", summary)

	return nil
}

func (r *Reader) stringToItem(data string) *pb.Item {
	slitted := strings.Split(data, ",")

	return &pb.Item{
		StndY:        r.convertToInt(slitted[0]),
		Sex:          r.convertToInt(slitted[1]),
		AgeGroup:     r.convertToInt(slitted[2]),
		Sido:         r.convertToInt(slitted[3]),
		Sgg:          r.convertToInt(slitted[4]),
		PersonID:     r.convertToInt(slitted[5]),
		KeySeq:       r.convertToInt(slitted[7]),
		YkihoID:      r.convertToInt(slitted[8]),
		RecuFrDt:     r.convertToInt(slitted[9]),
		DsbjtCd:      r.convertToInt(slitted[10]),
		DmdTramt:     r.convertToInt(slitted[11]),
		DmdSbrdnAmt:  r.convertToInt(slitted[12]),
		MainSick:     slitted[13],
		SubSick:      slitted[14],
		YkihoGubunCd: r.convertToInt(slitted[15]),
		YkihoSido:    r.convertToInt(slitted[16]),
	}
}

func (r *Reader) convertToInt(str string) int32 {
	i, _ := strconv.Atoi(str)

	return int32(i)
}

func (r *Reader) init() error {
	var err error

	r.conn, err = grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		return err
	}

	r.client = pb.NewHugeCSVClient(r.conn)

	return nil
}
