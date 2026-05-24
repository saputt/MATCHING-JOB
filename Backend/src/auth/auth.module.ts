import { Module } from "@nestjs/common";
import { JwtModule } from "@nestjs/jwt";
import { AuthController } from "./auth.controller";
import { AuthService } from "./auth.service";
import { AuthRepository } from "./auth.repository";
import { PrismaModule } from "../prisma/prisma.module";

@Module({
    providers:[AuthService, AuthRepository],
    controllers: [AuthController],
    imports: [
        PrismaModule,
        JwtModule.register({
            secret: process.env.SECRET_JWT,
            signOptions: { expiresIn : "7d" }
        })
    ]
})
export class AuthModule {}