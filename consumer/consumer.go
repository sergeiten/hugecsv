package consumer

import (
	"database/sql"
	"fmt"
	"io"
	"net"

	pb "github.com/sergeiten/hugecsv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Consumer struct {
	db *sql.DB
}

type Config struct {
	Host     string
	User     string
	Password string
	Database string
	Port     int
}

func New(config *Config) (*Consumer, error) {
	db, err := connect(config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		db: db,
	}, nil
}

func (consumer *Consumer) Send(stream pb.HugeCSV_SendServer) error {
	var processedCount int32

	for {
		item, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Summary{
				Processed: processedCount,
			})
		}
		if err != nil {
			fmt.Printf("failed to receieved stream %v", err)
		}
		processedCount++
		err = consumer.SavePatient(item)
		if err != nil {
			fmt.Printf("failed to save patient %v", err)
		}
		fmt.Printf("\rsaved %d", processedCount)
	}
}

func (consumer *Consumer) Serve() error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterHugeCSVServer(s, consumer)
	reflection.Register(s)

	return s.Serve(lis)
}

func connect(config *Config) (*sql.DB, error) {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", config.Host, config.User, config.Password, config.Port, config.Database)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (consumer *Consumer) SavePatient(h *pb.Item) error {
	query := `INSERT INTO Patients(
		stndy       
		,sex          
		,agegroup     
		,sido        
		,sgg          
		,personid     
		,keyseq      
		,ykihoid    
		,recufrdt  
		,dsbjtcd  
		,dmdtramt
		,dmdsbrdnamt
		,mainsick  
		,subsick  
		,ykihogubuncd
		,ykihosido  
	) VALUES (
		@stndy       
		,@sex          
		,@agegroup     
		,@sido        
		,@sgg          
		,@personid     
		,@keyseq      
		,@ykihoid    
		,@recufrdt  
		,@dsbjtcd  
		,@dmdtramt
		,@dmdsbrdnamt
		,@mainsick  
		,@subsick  
		,@ykihogubuncd
		,@ykihosido  
	)`

	_, err := consumer.db.Exec(query,
		sql.Named("stndy", h.StndY),
		sql.Named("sex", h.Sex),
		sql.Named("agegroup", h.AgeGroup),
		sql.Named("sido", h.Sido),
		sql.Named("sgg", h.Sgg),
		sql.Named("personid", h.PersonID),
		sql.Named("keyseq", h.KeySeq),
		sql.Named("ykihoid", h.YkihoID),
		sql.Named("recufrdt", h.RecuFrDt),
		sql.Named("dsbjtcd", h.DsbjtCd),
		sql.Named("dmdtramt", h.DmdTramt),
		sql.Named("dmdsbrdnamt", h.DmdSbrdnAmt),
		sql.Named("mainsick", h.MainSick),
		sql.Named("subsick", h.SubSick),
		sql.Named("ykihogubuncd", h.YkihoGubunCd),
		sql.Named("ykihosido", h.YkihoSido),
	)
	if err != nil {
		return err
	}

	return nil
}
