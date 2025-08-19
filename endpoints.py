from auth_service import AuthService
def register_user(username):
    service = AuthService()
    return service.create_user(username)
