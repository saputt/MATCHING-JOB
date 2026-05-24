import { Injectable } from "@nestjs/common";
import {PassportStrategy} from '@nestjs/passport'
import { ExtractJwt, Strategy } from "passport-jwt";

@Injectable()
export class JwtStrategy extends PassportStrategy(Strategy, "jwt") {
    constructor() {
        super({
            jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
            secretOrKey: process.env.SECRET_JWT!,
            ignoreExpiration: false
        })
    }
    
    async validate(payload : {id : string, email : string}) {
        const userData = {
            id : payload.id,
            email : payload.email
        }

        return userData
    }
}