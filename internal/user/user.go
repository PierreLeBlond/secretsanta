package user

type User struct {
  Name string
  Always []string
  Previous []string
}

func GetUserByName(users []*User, name string) *User {
  for _, user := range users {
    if user.Name == name {
      return user;
    }
  }
  return nil;
}
