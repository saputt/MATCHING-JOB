import { Injectable } from "@nestjs/common";
import { UserRepository } from "./user.repository";
import { UserProfileResponse } from "./interfaces/user-response.interface";
import { UpdateSkillsDto } from "./dto/update-skills.dto";

@Injectable()
export class UserService {
    constructor(private repo : UserRepository) {}

    async UpdateSkills(userData : UserProfileResponse, req : UpdateSkillsDto) {
        const oldSkills = userData.skills || []
        const newSkills = req.skills || []

        const uniqueSkillSet = new Set([...oldSkills, ...newSkills])
        const finalSkill = [...uniqueSkillSet]

        return this.repo.UpdateSkills(userData.id, finalSkill)
    }
}