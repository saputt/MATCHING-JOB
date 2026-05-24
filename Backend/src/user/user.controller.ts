import { Body, Controller, Param, Patch, Post } from "@nestjs/common";
import { UserService } from "./user.service";
import { UserExist } from "./pipes/user-exist.pipe";
import type { UserProfileResponse } from "./interfaces/user-response.interface";
import { UpdateSkillsDto } from "./dto/update-skills.dto";

@Controller("users")
export class UserController {
    constructor(private readonly service : UserService) {}

    @Patch('skills/:userId')
    async UpdateSkills(@Param("userId", UserExist) userData : UserProfileResponse, @Body() req : UpdateSkillsDto) {
        const res = await this.service.UpdateSkills(userData, req)
        return {
            message : "Update skills success",
            data : res
        }
    }
}