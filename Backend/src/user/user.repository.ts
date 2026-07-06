import { Injectable } from "@nestjs/common";
import { PrismaService } from "src/prisma/prisma.service";

@Injectable()
export class UserRepository {
    constructor(private prisma : PrismaService) {}

    async UpdateSkillsAndEmbedding(userId : string, skillsPostgresArray : string, vectorString: string) {
        return await this.prisma.$queryRawUnsafe<any[]>(`
            UPDATE "users"
            SET 
                skills = '${skillsPostgresArray}'::text[],
                embedding = '${vectorString}'::vector
            WHERE id = '${userId}'
            RETURNING skills, embedding::text AS embedding;;
        `);
    }

    async getUserEmbeddingString(userId: string) {
        return this.prisma.$queryRawUnsafe<any[]>(`
            SELECT embedding::text AS "embeddingStr"
            FROM "users"
            WHERE id = '${userId}'
            LIMIT 1;
        `);
    }

    async FindUserById(userId : string) {
        return this.prisma.user.findFirst({
            where : {
                id : userId
            }
        })
    }
}