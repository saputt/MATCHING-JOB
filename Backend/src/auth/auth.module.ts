import { Module } from "@nestjs/common";
import { JwtModule } from "@nestjs/jwt";
import { AuthController } from "./auth.controller";
import { AuthService } from "./auth.service";
import { AuthRepository } from "./auth.repository";
import { PrismaModule } from "../prisma/prisma.module";
import { JwtStrategy } from "./strategy/jwt-strategy";
import { PassportModule } from "@nestjs/passport";

@Module({
    providers:[AuthService, AuthRepository, JwtStrategy],
    controllers: [AuthController],
    imports: [
        PrismaModule,
        PassportModule.register({ defaultStrategy: 'jwt' }),
        JwtModule.register({
            secret: process.env.SECRET_JWT,
            signOptions: { expiresIn : "7d" }
        })
    ],
    exports : [PassportModule, JwtStrategy, AuthService]
})
export class AuthModule {}