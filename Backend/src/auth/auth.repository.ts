import { Injectable } from "@nestjs/common";
import { RegisterDto } from "./dto/register.dto";
import { PrismaService } from "../prisma/prisma.service";

@Injectable()
export class AuthRepository {
    constructor(private prisma : PrismaService) {}

    async createUser(dto : RegisterDto) {
        return this.prisma.user.create({
            data : {
                username : dto.username,
                email : dto.email,
                password : dto.password,
                skills : [],
            },
            select : {
                id : true,
                username : true,
                email : true,   
                createdAt : true             
            }
        })
    }

    async findUserByEmail(email : string) {
        return this.prisma.user.findUnique({
            where : {email}
        })
    }
}