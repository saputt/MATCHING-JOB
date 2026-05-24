import { Injectable } from "@nestjs/common";
import { PrismaService } from "src/prisma/prisma.service";

@Injectable()
export class UserRepository {
    constructor(private prisma : PrismaService) {}

    async UpdateSkills(userId : string, newSkill : string[]) {
        return this.prisma.user.update({
            where : {
                id : userId
            },
            data : {
                skills : newSkill
            },
            select : {
                skills: true
            }
        })
    }

    async FindUserById(userId : string) {
        return this.prisma.user.findFirst({
            where : {
                id : userId
            },
            select : {
                id : true,
                username : true,
                email : true,
                skills : true,
            }
        })
    }
}