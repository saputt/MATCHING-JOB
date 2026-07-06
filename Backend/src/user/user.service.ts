import { BadRequestException, Injectable, NotFoundException } from "@nestjs/common";
import { UserRepository } from "./user.repository";
import { AiServiceClient, GenerateEmbedding } from "src/ai/ai-service.client";
import { UpdatePreferencesDto } from "./dto/update-preferences.dto";

@Injectable()
export class UserService {
    constructor(
        private repo : UserRepository,
        private client : AiServiceClient
    ) {}

    async isUserExist(userId : string) {
        const userExist = await this.repo.FindUserById(userId)
        if (!userExist) throw new NotFoundException("user doesnt exist")
        return userExist
    }

    async getUserEmbedding(userId : string) {
        const embedding = await this.repo.getUserEmbeddingString(userId)
        if (!embedding || embedding.length == 0) throw new NotFoundException("user doesnt have embedding")
        return embedding[0].embeddingStr
    }

    async UpdatePreferences(userId : string, dto : UpdatePreferencesDto) {
        const userExist = await this.isUserExist(userId)
        const oldSkills = userExist?.skills || []
        const newSkills = dto.skills || []

        const uniqueSkillSet = new Set([...oldSkills, ...newSkills])
        const finalSkill = [...uniqueSkillSet] 

        if (userExist.titleInterest == "" && dto.titleInterest == "") throw new BadRequestException(" title interese cannot be empty")

        const payloadEmbedding : GenerateEmbedding = {
            hardSkills : finalSkill,
            title : dto.titleInterest ?? userExist.titleInterest
        }

        const textEmbedding = await this.client.generateEmbedding(payloadEmbedding)

        const skillsPostgresArray = `{${finalSkill.map(s => `"${s}"`).join(',')}}`
        const vectorString = `[${textEmbedding.data.join(',')}]`;

        return this.repo.UpdateSkillsAndEmbedding(userId, skillsPostgresArray, vectorString)
    }
}