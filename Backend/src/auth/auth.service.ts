import { ForbiddenException, Injectable, NotFoundException, UnauthorizedException } from "@nestjs/common";
import { AuthRepository } from "./auth.repository";
import { LoginDto } from "./dto/login.dto";
import { JwtService } from "@nestjs/jwt";
import { RegisterDto } from "./dto/register.dto";
import { hashing } from "src/common/helpers/hash.helper";

@Injectable()
export class AuthService {
    constructor(
        private repo : AuthRepository,
        private jwt : JwtService
    ) 
    {}

    async signToken(id : string, email : string) {
        const payload = {id, email}

        const token = await this.jwt.signAsync(payload, {
            secret : process.env.SECRET_JWT,
            expiresIn : "7d"
        })

        return token
    }

    async login(dto : LoginDto) {
        const isUserExist = await this.repo.findUserByEmail(dto.email)

        if(!isUserExist) throw new UnauthorizedException("Invalid credential")

        const isPasswordCorrect = await hashing.compare(dto.password, isUserExist.password)

        if (!isPasswordCorrect) throw new ForbiddenException("Invalid credential")

        const token = await this.signToken(
            isUserExist.id,
            isUserExist.email
        )

        return {
            user : {
                id : isUserExist.id,
                email : isUserExist.email,
            },
            token
        }
    }

    async register(dto : RegisterDto) {
        const isEmailExist = await this.repo.findUserByEmail(dto.email)

        if (isEmailExist) throw new ForbiddenException("Email already exist")

        const userData = {
            username : dto.username,
            email : dto.email,
            password : await hashing.hash(dto.password)
        }

        return this.repo.createUser(userData)
    }
}