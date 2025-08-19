from models import User
class AuthService:
    def create_user(self, username):
        return User(username)
    def validate_username(self, username):
        return len(username) >= 3
    def create_user(self, username):
        if not self.validate_username(username):
            raise ValueError('Username too short')
        return User(username)
