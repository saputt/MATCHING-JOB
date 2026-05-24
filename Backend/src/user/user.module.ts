import { Module } from "@nestjs/common";
import { UserService } from "./user.service";
import { UserRepository } from "./user.repository";
import { UserController } from "./user.controller";
import { PrismaModule } from "src/prisma/prisma.module";

@Module({
    providers : [UserService, UserRepository],
    controllers : [UserController],
    imports : [PrismaModule]
})
export class UserModule {}