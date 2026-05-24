import { Body, Controller, Post } from "@nestjs/common";
import { AuthService } from "./auth.service";
import { LoginDto } from "./dto/login.dto";
import { RegisterDto } from "./dto/register.dto";

@Controller("auth")
export class AuthController {
    constructor (private readonly service : AuthService) {}

    @Post("login")
    async login(@Body() req : LoginDto) {
        const res = await this.service.login(req)
        return {
            message : "Login success",
            data : res
        }
    }

    @Post("register")
    async register(@Body() req : RegisterDto) {
        const res = await this.service.register(req)
        return {
            message : "Register success",
            data : res
        }
    }
}