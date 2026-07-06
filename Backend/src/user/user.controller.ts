import { Body, Controller, Param, Patch, Post, UseGuards } from "@nestjs/common";
import { UserService } from "./user.service";
import type { UserProfileResponse } from "./interfaces/user-response.interface";
import { GetUser } from "src/common/decorators/get-user.decorator";
import { UpdatePreferencesDto } from "./dto/update-preferences.dto";
import { JwtAuthGuard } from "src/common/guards/jwt-auth.guard";

@Controller("user")
@UseGuards(JwtAuthGuard)
export class UserController {
    constructor(private readonly service : UserService) {}

    @Patch('preferences')
    async UpdateSkills(@GetUser('id') userId : string, @Body() dto : UpdatePreferencesDto) {
        const res = await this.service.UpdatePreferences(userId, dto)
        return {
            message : "Update preferences success",
            data : res
        }
    }
}