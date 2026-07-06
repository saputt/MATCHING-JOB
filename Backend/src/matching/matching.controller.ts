import { Controller, Get, UseGuards } from "@nestjs/common";
import { MatchingService } from "./matching.service";
import { GetUser } from "src/common/decorators/get-user.decorator";
import { JwtAuthGuard } from "src/common/guards/jwt-auth.guard";

@Controller("matching")
@UseGuards(JwtAuthGuard)
export class MatchingController {
    constructor(private readonly service : MatchingService) {}

    @Get("recommendation")
    async getRecommendation(@GetUser('id') userId : string) {
        const res = await this.service.getTopRecomendation(userId)
        return {
            message : "get top 10 match job success",
            data : res
        }
    }
}