package model

import "testing"

func TestRelationshipValidate(t *testing.T) {
	tests := []struct {
		name string
		rel  Relationship
		want string
	}{
		{
			name: "valid relationship",
			rel:  Relationship{FollowerID: 1, FollowedID: 2},
		},
		{
			name: "missing follower id",
			rel:  Relationship{FollowedID: 2},
			want: "follower_id can't be blank",
		},
		{
			name: "missing followed id",
			rel:  Relationship{FollowerID: 1},
			want: "followed_id can't be blank",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rel.Validate()
			if tt.want == "" {
				if err != nil {
					t.Fatalf("Validate() returned unexpected error: %v", err)
				}
				return
			}
			if err == nil || err.Error() != tt.want {
				t.Fatalf("Validate() error = %v, want %q", err, tt.want)
			}
		})
	}
}
