package colorextra

import "testing"

func TestConvertUint8ToUint32Color(t *testing.T) {
	type args struct {
		in uint8
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			"0",
			args{in: 0},
			0,
		}, {
			"255",
			args{in: 255},
			65535,
		}, {
			"127",
			args{in: 127},
			32639,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertUint8ToUint32Color(tt.args.in); got != tt.want {
				t.Errorf("ConvertUint8ToUint32Color() = %v, want %v", got, tt.want)
			}
		})
	}
}
